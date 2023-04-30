use std::time::Duration;

use indicatif::{ProgressBar, ProgressStyle};

pub fn run<R, F: FnMut() -> R>(
    mut action: F,
    pending_message: String,
    finished_message: String,
) -> R {
    let spinner = ProgressBar::new_spinner();
    spinner.enable_steady_tick(Duration::from_millis(100));
    spinner.set_style(
        ProgressStyle::default_spinner()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"])
            .template("{spinner:.green} {msg:.green}")
            .unwrap(),
    );
    spinner.set_message(pending_message);
    let result = action();
    spinner.finish_with_message(finished_message);
    result
}
