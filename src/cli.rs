use clap::Parser;
use keylight_on::keylight::CommandLight;

#[derive(Parser, Debug)]
#[command(name = "keylight-on", about = "Turn on your keylight")]
struct Args {
    #[clap(
        short,
        long,
        value_name = "index",
        help = "Index in the list of lights that will be changed"
    )]
    index: usize,
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
        index: args.index,
        on: args.on,
        brightness: args.brightness,
        temperature: args.temperature,
    }
}
