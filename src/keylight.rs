pub enum KeylightError {
    CommandError(String),
}

#[derive(Debug)]
pub struct Metadata {
    name: String,
    ip: String,
    port: u16,
}

impl PartialEq for Metadata {
    fn eq(&self, other: &Self) -> bool {
        self.name == other.name && self.ip == other.ip && self.port == other.port
    }
}

pub struct OperatingData {
    on: bool,
    brightness: i32,
    temperature: i32,
}
pub struct Keylight {
    metadata: Metadata,
    operating_data: OperatingData,
}

#[derive(Debug)]
pub struct DiscoveredKeylight {
    metadata: Metadata,
}

impl PartialEq for DiscoveredKeylight {
    fn eq(&self, other: &Self) -> bool {
        self.metadata == other.metadata
    }
}

impl DiscoveredKeylight {
    pub fn new(name: String, ip: String, port: u16) -> DiscoveredKeylight {
        DiscoveredKeylight {
            metadata: Metadata { name, ip, port },
        }
    }

    pub fn connect(self, operating_data: OperatingData) -> Keylight {
        Keylight {
            metadata: self.metadata,
            operating_data,
        }
    }
}

impl Keylight {
    pub fn update(&mut self, operating_data: OperatingData) -> &Self {
        unimplemented!()
    }

    pub fn turn_on(&mut self) -> Result<&Self, KeylightError> {
        unimplemented!()
    }

    pub fn change_brightness(&mut self, brightness: i32) -> Result<&Self, KeylightError> {
        unimplemented!()
    }

    pub fn change_temperature(&mut self, temperature: i32) -> Result<&Self, KeylightError> {
        unimplemented!()
    }

    pub fn change(
        &mut self,
        on: bool,
        brightness: i32,
        temperature: i32,
    ) -> Result<&Self, KeylightError> {
        unimplemented!()
    }
}
