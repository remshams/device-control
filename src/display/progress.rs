use std::time::Duration;

use indicatif::{ProgressBar, ProgressStyle};

pub fn run<F: FnMut() -> ()>(mut action: F, pending_message: String, finished_message: String) {
    let spinner = ProgressBar::new_spinner();
    spinner.enable_steady_tick(Duration::from_millis(100));
    spinner.set_style(
        ProgressStyle::default_spinner()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"])
            .template("{spinner:.white} {msg}")
            .unwrap(),
    );
    spinner.set_message(pending_message);
    action();
    spinner.finish_with_message(finished_message);
}
