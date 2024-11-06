use prost::Message;
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::util::Timeout;
use rdkafka::ClientConfig;
use schema_registry_converter::async_impl::easy_proto_raw::EasyProtoRawEncoder;
use schema_registry_converter::async_impl::schema_registry::SrSettings;
use schema_registry_converter::schema_registry_common::SubjectNameStrategy;
use socket2::{Domain, Socket, Type};
use std::error::Error;
use std::net::SocketAddr;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::net::UdpSocket;
use tokio::sync::mpsc;

pub mod iot {
    include!(concat!(env!("OUT_DIR"), "/production.iot.rs"));
}
use iot::Motor;

struct ProcessingContext {
    producer: FutureProducer,
    encoder: EasyProtoRawEncoder,
    topic: String,
    message_count: AtomicUsize,
    start_time: Instant,
}

fn create_producer(bootstrap_server: &str) -> FutureProducer {
    ClientConfig::new()
        .set("bootstrap.servers", bootstrap_server)
        .set("batch.size", "20971520") // 20 MiB
        .set("linger.ms", "2")
        .set("compression.type", "zstd")
        .create()
        .expect("Failed to create producer")
}

async fn process_message(
    context: Arc<ProcessingContext>,
    motor: Motor,
) -> Result<(), Box<dyn Error + Send + Sync>> {
    let mut buf = Vec::new();
    motor.encode(&mut buf)?;

    let proto_payload = context
        .encoder
        .encode(
            &buf,
            "production.iot.Motor",
            SubjectNameStrategy::TopicNameStrategy(context.topic.clone(), false),
        )
        .await
        .map_err(|e| Box::new(e) as Box<dyn Error + Send + Sync>)?;

    let motor_id = motor
        .id
        .clone()
        .unwrap_or_else(|| "unknown_motor".to_string());
    context
        .producer
        .send(
            FutureRecord::to(&context.topic)
                .payload(&proto_payload)
                .key(motor_id.as_bytes()),
            Timeout::After(Duration::from_secs(1)),
        )
        .await
        .map_err(|(err, _)| Box::new(err) as Box<dyn Error + Send + Sync>)?;

    let count = context.message_count.fetch_add(1, Ordering::Relaxed) + 1;
    if count % 100000 == 0 {
        let duration = context.start_time.elapsed();
        println!(
            "Total messages sent to Kafka: {}, Time elapsed: {:?}, Avg msg/sec: {:.2}",
            count,
            duration,
            count as f64 / duration.as_secs_f64()
        );
    }
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let start_time = Instant::now();
    println!("Starting UDP listener and Kafka producer...");

    // Set up UDP socket
    let socket = Socket::new(Domain::IPV4, Type::DGRAM, None)?;
    socket.set_reuse_address(true)?;
    socket.set_recv_buffer_size(20 * 1024 * 1024)?;
    let addr: SocketAddr = "127.0.0.1:50537".parse()?;
    socket.bind(&addr.into())?;
    let socket = UdpSocket::from_std(socket.into())?;
    println!("Listening on 127.0.0.1:50537");

    // Set up Kafka producer and Schema Registry encoder
    let bootstrap_server = "localhost:9092";
    let producer = create_producer(bootstrap_server);
    let schema_registry_url = "https://confluent-schema-registry.internal.hongbomiao.com";
    let sr_settings = SrSettings::new(schema_registry_url.to_string());
    let encoder = EasyProtoRawEncoder::new(sr_settings);
    let topic = "production.iot.motor.proto".to_string();

    // Create processing context with start_time
    let context = Arc::new(ProcessingContext {
        producer,
        encoder,
        topic,
        message_count: AtomicUsize::new(0),
        start_time,
    });

    // Create channel for message processing
    let (tx, mut rx) = mpsc::channel::<Motor>(1000000); // Buffer size of 1000000 messages

    // Spawn message processing task
    let processing_context = context.clone();
    tokio::spawn(async move {
        let mut handles = Vec::new();
        while let Some(motor) = rx.recv().await {
            let ctx = processing_context.clone();
            handles.push(tokio::spawn(async move {
                if let Err(e) = process_message(ctx, motor).await {
                    eprintln!("Error processing message: {}", e);
                }
            }));
            handles.retain(|handle| !handle.is_finished());
            while handles.len() >= 2000 {
                futures::future::join_all(handles.drain(..1000)).await;
            }
        }
    });

    // Main UDP receiving loop
    let mut buffer = vec![0u8; 65536];
    loop {
        let (amt, _src) = socket.recv_from(&mut buffer).await?;
        match Motor::decode(&buffer[..amt]) {
            Ok(motor) => {
                if let Err(e) = tx.send(motor).await {
                    eprintln!("Failed to send message to processing task: {}", e);
                }
            }
            Err(e) => {
                eprintln!("Failed to decode protobuf message: {}", e);
            }
        }
    }
}
