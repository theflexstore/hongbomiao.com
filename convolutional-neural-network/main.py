import wandb
import yaml
from evaluate import evaluate
from model.data_loader import train_data_loader, val_data_loader
from model.net import net
from torch import nn, optim
from train import train
from utils.device import device
from utils.writer import write_params


def main():
    with open("params.yaml", "r") as f:
        params = yaml.safe_load(f)

    write_params(params)

    with wandb.init(
        entity="hongbo-miao", project="convolutional-neural-network", config=params
    ) as wb:
        config = wb.config
        net.to(device)

        criterion = nn.CrossEntropyLoss()
        optimizer = optim.SGD(net.parameters(), lr=config["lr"], momentum=0.9)
        # optimizer = optim.Adam(net.parameters(), lr=config.learning_rate)

        for epoch in range(config["train"]["epochs"]):
            train_loss = train(net, train_data_loader, optimizer, criterion)
            train_acc = evaluate(net, train_data_loader)
            val_acc = evaluate(net, val_data_loader)

            print({"Train": train_acc, "Validation": val_acc})
            wb.log(
                {
                    "epoch": epoch,
                    "train_loss": train_loss,
                    "train_acc": train_acc,
                    "val_acc": val_acc,
                }
            )


if __name__ == "__main__":
    main()
