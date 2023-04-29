use clap::{Parser, Subcommand};
use keylight_on::keylight::{KeylightCommand, LightCommand};

#[derive(Parser, Debug)]
#[command(name = "keylight-on", about = "Turn on your keylight")]
struct Args {
    #[clap(subcommand)]
    keylight_command: Command,
}

#[derive(Subcommand, Debug)]
enum Command {
    #[command(name = "list", about = "List available keylights")]
    List,
    #[command(name = "send-command", about = "Set keylight parameters")]
    SendCommand(CommandLightArgs),
}

#[derive(Parser, Debug)]
struct CommandLightArgs {
    #[clap(short, long, value_name = "id", help = "Id of controlled keylight")]
    id: String,
    #[clap(
        short,
        long,
        value_name = "index",
        help = "Index in the list of lights that will be changed"
    )]
    #[clap(
        short,
        long,
        value_name = "light_index",
        help = "Index of keylight light, in most cases 0. Defaults to 0 if not provided"
    )]
    light_index: Option<usize>,
    #[clap(
        short,
        long,
        value_name = "on",
        help = "Turn light on or off",
        required = false
    )]
    on: Option<bool>,
    #[clap(short, long, value_name = "brightness", required = false)]
    brightness: Option<u16>,
    #[clap(short, long, value_name = "temperature", required = false)]
    temperature: Option<u16>,
}

pub fn parse() -> KeylightCommand {
    let args: Args = Args::parse();
    match args.keylight_command {
        Command::SendCommand(command_light_args) => KeylightCommand::SendCommand(LightCommand {
            id: command_light_args.id,
            index: command_light_args.light_index.unwrap_or(0),
            on: command_light_args.on,
            brightness: command_light_args.brightness,
            temperature: command_light_args.temperature,
        }),
        Command::List => KeylightCommand::List,
    }
}
