use clap::Parser;
use keylight_on::keylight::CommandLight;

#[derive(Parser, Debug)]
#[command(name = "keylight-on", about = "Turn on your keylight")]
struct Args {
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

pub fn parse() -> CommandLight {
    let args: Args = Args::parse();
    CommandLight {
        id: args.id,
        index: args.light_index.unwrap_or(0),
        on: args.on,
        brightness: args.brightness,
        temperature: args.temperature,
    }
}
