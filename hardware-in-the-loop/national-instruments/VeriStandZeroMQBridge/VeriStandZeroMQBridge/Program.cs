using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Google.Protobuf;
using NationalInstruments.VeriStand.ClientAPI;
using NetMQ;
using NetMQ.Sockets;
using VeriStandZeroMQBridge;

public class Configuration
{
    public string VeristandGatewayHost { get; set; }
    public string SystemDefinitionPath { get; set; }
    public uint VeriStandConnectionTimeoutMs { get; set; }
    public int ZeroMQPort { get; set; }
    public string ZeroMQAddress => $"tcp://*:{ZeroMQPort}";
    public int ProducerNumber { get; set; }
    public int TargetFrequencyHz { get; set; }
    public int CalibrationTimeS { get; set; }

    public static Configuration Load()
    {
        string env = Environment.GetEnvironmentVariable("APP_ENVIRONMENT") ?? "Development";
        string envFile = $".env.{env.ToLower()}";

        if (!File.Exists(envFile))
        {
            throw new FileNotFoundException($"Environment file not found: {envFile}");
        }

        DotNetEnv.Env.Load(envFile);

        return new Configuration
        {
            VeristandGatewayHost = Environment.GetEnvironmentVariable("VERISTAND_GATEWAY_HOST"),
            SystemDefinitionPath = Environment.GetEnvironmentVariable("SYSTEM_DEFINITION_PATH"),
            VeriStandConnectionTimeoutMs = uint.Parse(
                Environment.GetEnvironmentVariable("VERISTAND_CONNECTION_TIMEOUT_MS")
            ), // Changed to uint.Parse
            ZeroMQPort = int.Parse(Environment.GetEnvironmentVariable("ZEROMQ_PORT")),
            ProducerNumber = int.Parse(Environment.GetEnvironmentVariable("PRODUCER_NUMBER")),
            TargetFrequencyHz = int.Parse(
                Environment.GetEnvironmentVariable("TARGET_FREQUENCY_HZ")
            ),
            CalibrationTimeS = int.Parse(Environment.GetEnvironmentVariable("CALIBRATION_TIME_S")),
        };
    }
}

public record DataQualityMetrics
{
    public int TotalReadAttempts { get; set; }
    public int FailedReads { get; set; }
    public int LastGoodValueUsageCount { get; set; }
    public DateTime LastGoodValueTimestamp { get; set; }
    public TimeSpan LongestLastGoodValuePeriod { get; set; }
    public double LastGoodValuePercentage =>
        TotalReadAttempts > 0 ? (LastGoodValueUsageCount * 100.0 / TotalReadAttempts) : 0;

    public override string ToString()
    {
        return $"Total Reads: {TotalReadAttempts}, "
            + $"Failed: {FailedReads}, "
            + $"Last Good Value Usage: {LastGoodValuePercentage:F2}%, "
            + $"Longest Last Good Value Period: {LongestLastGoodValuePeriod.TotalMilliseconds:F0}ms";
    }
}

public static class StringSanitizer
{
    private static (int leading, int trailing) CountLeadingTrailingUnderscores(string s)
    {
        int leadingCount = s.Length - s.TrimStart('_').Length;
        int trailingCount = s.Length - s.TrimEnd('_').Length;
        return (leadingCount, trailingCount);
    }

    public static string SanitizeSignalName(string s)
    {
        if (string.IsNullOrEmpty(s))
            return s;

        (int leadingUnderscores, int trailingUnderscores) = CountLeadingTrailingUnderscores(s);

        // Replace special characters
        s = s.Replace("Δ", "d")
            .Replace("²", "2")
            .Replace("³", "3")
            .Replace("°", "deg")
            .Replace("℃", "degc")
            .Replace("%", "pct")
            .Replace("Ω", "ohm")
            .Replace("/", "."); // Convert forward slashes to dots

        // Replace special characters with underscore, except for dots
        s = Regex.Replace(s, @"[^a-zA-Z0-9_\.]", "_");

        // Replace multiple underscores with single underscore
        s = Regex.Replace(s, @"_+", "_");

        // Remove leading and trailing underscores
        s = s.Trim('_');

        // Add back original leading and trailing underscores
        s = new string('_', leadingUnderscores) + s + new string('_', trailingUnderscores);

        return s;
    }
}

public class Program
{
    private static Configuration config;
    private static Channel<byte[]> channel;
    private static volatile bool isRunning = true;
    private static int messageCount = 0;
    private static int messagesLastInterval = 0;
    private static readonly DataQualityMetrics metrics = new DataQualityMetrics();
    private static long nextGlobalTimestampNs =
        DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() * 1000;
    private static readonly object timestampLock = new object();

    private static long GetNextTimestamp()
    {
        lock (timestampLock)
        {
            long timestampNs = nextGlobalTimestampNs;
            nextGlobalTimestampNs += (1000000 / config.TargetFrequencyHz);
            return timestampNs * 1000;
        }
    }

    public static async Task Main(string[] args)
    {
        try
        {
            config = Configuration.Load();
            await Console.Out.WriteLineAsync(
                $"[Config] Loaded configuration for environment: {Environment.GetEnvironmentVariable("APP_ENVIRONMENT") ?? "Development"}"
            );

            channel = Channel.CreateUnbounded<byte[]>(
                new UnboundedChannelOptions { SingleReader = true, SingleWriter = false }
            );

            using (PublisherSocket publisher = new PublisherSocket())
            {
                publisher.Bind(config.ZeroMQAddress);
                await Console.Out.WriteLineAsync(
                    $"[ZMQ] Publisher bound to {config.ZeroMQAddress}"
                );

                try
                {
                    IWorkspace2 workspace = new Factory().GetIWorkspace2(
                        config.VeristandGatewayHost
                    );
                    string[] aliases,
                        channels;
                    workspace.GetAliasList(out aliases, out channels);
                    await Console.Out.WriteLineAsync(
                        $"[Config] Aliases Count: {aliases.Length} | Aliases: {string.Join(", ", aliases)}"
                    );
                    await Console.Out.WriteLineAsync(
                        $"[Config] Channels Count: {channels.Length} | Channels: {string.Join(", ", channels)}"
                    );

                    workspace.ConnectToSystem(
                        config.SystemDefinitionPath,
                        true,
                        config.VeriStandConnectionTimeoutMs
                    );
                    await Console.Out.WriteLineAsync("[Status] Data collection started");

                    double[] calibrationValues = new double[channels.Length];
                    double[] calibrationLastGoodValues = new double[channels.Length];
                    var calibrationResults = await PerformCalibration(
                        workspace,
                        channels,
                        calibrationValues,
                        calibrationLastGoodValues
                    );

                    Stopwatch totalStopwatch = Stopwatch.StartNew();

                    List<Task> producerTasks = new List<Task>();
                    for (int i = 0; i < config.ProducerNumber; i++)
                    {
                        producerTasks.Add(
                            Task.Run(() => ProducerTask(workspace, channels, calibrationResults))
                        );
                    }

                    Task consumerTask = Task.Run(async () =>
                    {
                        await foreach (byte[] data in channel.Reader.ReadAllAsync())
                        {
                            // Signals signals = Signals.Parser.ParseFrom(data);
                            // DateTime datetime = DateTimeOffset
                            //     .FromUnixTimeMilliseconds(signals.TimestampNs / 1000000)
                            //     .DateTime;
                            // Console.WriteLine($"Timestamp: {signals.TimestampNs} ({datetime})");

                            publisher.SendFrame(data);
                            Interlocked.Increment(ref messageCount);
                            Interlocked.Increment(ref messagesLastInterval);
                        }
                    });

                    Task monitoringTask = Task.Run(MonitorFrequency);

                    await Task.WhenAll(
                        producerTasks.Concat(new[] { consumerTask, monitoringTask })
                    );

                    double totalTime = totalStopwatch.Elapsed.TotalSeconds;
                    double averageMessagesPerSecond = messageCount / totalTime;
                    await Console.Out.WriteLineAsync(
                        $"[Complete] Runtime: {totalTime:F2}s | Avg Speed: {averageMessagesPerSecond:F2} msg/s | Total Messages: {messageCount:N0}"
                    );
                }
                catch (Exception ex)
                {
                    await Console.Out.WriteLineAsync($"[ERROR] {ex.Message}");
                    Environment.Exit(1);
                }
            }
        }
        catch (Exception ex)
        {
            await Console.Out.WriteLineAsync($"[ERROR] Configuration error: {ex.Message}");
            Environment.Exit(1);
        }
    }

    // ... rest of your existing methods, replacing const references with config properties ...
    private static Signals CreateSignalsMessage(
        string[] channels,
        double[] values,
        bool useLastGoodValues,
        double[] lastGoodValues
    )
    {
        return new Signals
        {
            TimestampNs = GetNextTimestamp(),
            Signals_ =
            {
                channels.Zip(
                    useLastGoodValues ? lastGoodValues : values,
                    (name, value) =>
                        new Signal
                        {
                            Name = StringSanitizer.SanitizeSignalName(name),
                            Value = (float)value,
                            IsLastGoodValue = useLastGoodValues,
                        }
                ),
            },
            SkippedTickNumber = 0,
            IsUsingLastGoodValues = useLastGoodValues,
            FrequencyHz = config.TargetFrequencyHz,
        };
    }

    private static async Task<(int readEveryNSampleCount, long periodTicks)> PerformCalibration(
        IWorkspace2 workspace,
        string[] channels,
        double[] values,
        double[] lastGoodValues
    )
    {
        List<long> sampleTimes = new List<long>();
        Stopwatch calibrationStopwatch = Stopwatch.StartNew();
        long calibrationEndTime = config.CalibrationTimeS * Stopwatch.Frequency;

        Stopwatch stopwatch = new Stopwatch();
        stopwatch.Start();
        long periodTicks = Stopwatch.Frequency / config.TargetFrequencyHz;
        long nextTick = stopwatch.ElapsedTicks;

        Console.WriteLine("Starting calibration...");

        while (calibrationStopwatch.ElapsedTicks < calibrationEndTime && isRunning)
        {
            if (stopwatch.ElapsedTicks >= nextTick)
            {
                long beforeRead = calibrationStopwatch.ElapsedTicks;

                try
                {
                    workspace.GetMultipleChannelValues(channels, out values);
                    long readTime = calibrationStopwatch.ElapsedTicks - beforeRead;
                    sampleTimes.Add(readTime);
                    Array.Copy(values, lastGoodValues, values.Length);
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Calibration error: {ex.Message}");
                }

                nextTick += periodTicks;
                HandlePreciseTiming(nextTick);
            }
        }

        return CalculateCalibrationResults(sampleTimes);
    }

    private static (int readEveryNSampleCount, long periodTicks) CalculateCalibrationResults(
        List<long> sampleTimes
    )
    {
        long averageReadTime = (long)sampleTimes.Average();
        long minReadTime = sampleTimes.Min();
        long maxReadTime = sampleTimes.Max();
        double stdDev = Math.Sqrt(sampleTimes.Average(v => Math.Pow(v - averageReadTime, 2)));
        int sampleCount = sampleTimes.Count;

        long[] sortedTimes = sampleTimes.OrderBy(x => x).ToArray();
        long percentile95 = sortedTimes[(int)(sortedTimes.Length * 0.95)];
        long minReadPeriodTicks = percentile95 * 2;
        int actualReadFrequency = (int)(Stopwatch.Frequency / minReadPeriodTicks);
        int readEveryNSampleCount = Math.Max(1, config.TargetFrequencyHz / actualReadFrequency);

        Console.WriteLine("\n========== CALIBRATION RESULTS ==========");
        Console.WriteLine($"Samples collected: {sampleCount}");
        Console.WriteLine(
            $"Average read time: {averageReadTime * 1000.0 / Stopwatch.Frequency:F3}ms"
        );
        Console.WriteLine($"Min read time: {minReadTime * 1000.0 / Stopwatch.Frequency:F3}ms");
        Console.WriteLine($"Max read time: {maxReadTime * 1000.0 / Stopwatch.Frequency:F3}ms");
        Console.WriteLine($"Standard deviation: {stdDev * 1000.0 / Stopwatch.Frequency:F3}ms");
        Console.WriteLine($"Measured frequency: {Stopwatch.Frequency / averageReadTime:F1} Hz");
        Console.WriteLine($"Safe read frequency: {actualReadFrequency} Hz");
        Console.WriteLine($"Reading every {readEveryNSampleCount} samples");
        Console.WriteLine("========================================\n");

        return (readEveryNSampleCount, Stopwatch.Frequency / config.TargetFrequencyHz);
    }

    private static void HandlePreciseTiming(long nextTick)
    {
        long sleepTicks = nextTick - Stopwatch.GetTimestamp();
        if (sleepTicks > 0)
        {
            long sleepMillis = sleepTicks * 1000 / Stopwatch.Frequency;
            if (sleepMillis > 1)
            {
                Task.Delay(1).Wait();
            }
            SpinWait.SpinUntil(() => Stopwatch.GetTimestamp() >= nextTick);
        }
    }

    private static async Task ProducerTask(
        IWorkspace2 workspace,
        string[] channels,
        (int readEveryNSampleCount, long periodTicks) calibrationResults
    )
    {
        double[] values = new double[channels.Length];
        double[] lastGoodValues = new double[channels.Length];

        var (readEveryNSampleCount, periodTicks) = calibrationResults;

        Stopwatch stopwatch = new Stopwatch();
        stopwatch.Start();
        long nextTick = stopwatch.ElapsedTicks;

        while (isRunning)
        {
            if (stopwatch.ElapsedTicks >= nextTick)
            {
                bool useLastGoodValues = false;
                try
                {
                    if (metrics.TotalReadAttempts % readEveryNSampleCount == 0)
                    {
                        workspace.GetMultipleChannelValues(channels, out values);
                        Array.Copy(values, lastGoodValues, values.Length);
                        metrics.LastGoodValueTimestamp = DateTime.UtcNow;
                    }
                    else
                    {
                        useLastGoodValues = true;
                    }

                    UpdateMetrics(useLastGoodValues);

                    Signals signals = CreateSignalsMessage(
                        channels,
                        values,
                        useLastGoodValues,
                        lastGoodValues
                    );
                    byte[] data = signals.ToByteArray();
                    await channel.Writer.WriteAsync(data);

                    nextTick += periodTicks;
                    HandlePreciseTiming(nextTick);
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error during normal operation: {ex.Message}");
                    metrics.FailedReads++;
                    await Task.Delay(10);
                }
            }
        }
    }

    private static void UpdateMetrics(bool useLastGoodValues)
    {
        metrics.TotalReadAttempts++;
        if (useLastGoodValues)
        {
            metrics.LastGoodValueUsageCount++;
            TimeSpan currentLgvPeriod = DateTime.UtcNow - metrics.LastGoodValueTimestamp;
            if (currentLgvPeriod > metrics.LongestLastGoodValuePeriod)
            {
                metrics.LongestLastGoodValuePeriod = currentLgvPeriod;
            }
        }
    }

    private static void MonitorFrequency()
    {
        Stopwatch stopwatch = new Stopwatch();
        stopwatch.Start();
        int lastCount = 0;
        long lastTime = 0;

        while (isRunning)
        {
            Thread.Sleep(1000);

            long currentTime = stopwatch.ElapsedMilliseconds;
            int currentCount = messageCount;
            int messages = currentCount - lastCount;
            double actualFrequency = messages / ((currentTime - lastTime) / 1000.0);

            Console.WriteLine($"Actual frequency: {actualFrequency:F2} Hz");
            Console.WriteLine($"Target frequency: {config.TargetFrequencyHz} Hz");
            Console.WriteLine($"Difference: {actualFrequency - config.TargetFrequencyHz:F2} Hz");
            Console.WriteLine($"Data Quality: {metrics}");

            lastCount = currentCount;
            lastTime = currentTime;
        }
    }
}
