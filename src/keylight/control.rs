use super::{model::KeylightError, Keylight, KeylightMetadata, Light};
use log::debug;

pub trait KeylightFinder {
    type Output: IntoIterator<Item = KeylightMetadata>;

    fn discover(&self) -> Self::Output;
}

pub trait KeylightAdapter {
    fn lights(&self, ip: &str, port: &u16) -> Result<Vec<Light>, KeylightError>;
    fn set_lights(&self, ip: &str, lights: &[Light]) -> Result<(), KeylightError>;
}

pub trait KeylightDb {
    fn store(&self, metadatas: &[&KeylightMetadata]) -> Result<(), KeylightError>;
    fn load(&self) -> Result<Vec<KeylightMetadata>, KeylightError>;
}

pub struct KeylightControl<'a, F: KeylightFinder, A: KeylightAdapter, Db: KeylightDb> {
    keylight_finder: &'a F,
    keylight_adapter: &'a A,
    keylight_db: &'a Db,
    pub lights: Vec<Keylight<'a, A>>,
}

impl<'a, F: KeylightFinder, A: KeylightAdapter, Db: KeylightDb> KeylightControl<'a, F, A, Db> {
    pub fn new(
        keylight_finder: &'a F,
        keylight_adapter: &'a A,
        keylight_db: &'a Db,
    ) -> KeylightControl<'a, F, A, Db> {
        KeylightControl {
            keylight_finder,
            keylight_adapter,
            keylight_db,
            lights: vec![],
        }
    }

    pub fn find_keylight_mut(&mut self, id: &str) -> Option<&mut Keylight<'a, A>> {
        self.lights.iter_mut().find(|light| light.metadata.id == id)
    }

    pub fn find_keylight(&self, id: &str) -> Option<&Keylight<'a, A>> {
        self.lights.iter().find(|light| light.metadata.id == id)
    }

    pub fn load_keylights(&mut self) -> Result<(), KeylightError> {
        let result = self.keylight_db.load();
        match result {
            Ok(metadata) => {
                self.lights = metadata
                    .into_iter()
                    .map(|metadata| Keylight::new(self.keylight_adapter, metadata, None))
                    .collect();
                debug!("Loaded {} keylights", self.lights.len());

                Ok(())
            }
            Err(_) => {
                self.discover_keylights();
                self.store_keylights()
            }
        }
    }

    pub fn discover_keylights(&mut self) {
        self.lights = self
            .keylight_finder
            .discover()
            .into_iter()
            .map(|metadata| Keylight::new(self.keylight_adapter, metadata, None))
            .collect();
        self.deduplicate_keylights();
        debug!("Discovered {} keylights", self.lights.len());
    }

    pub fn store_keylights(&self) -> Result<(), KeylightError> {
        let keylight_metadatas = self
            .lights
            .iter()
            .map(|light| &light.metadata)
            .collect::<Vec<&KeylightMetadata>>();
        debug!("Storing {} keylights", keylight_metadatas.len());
        self.keylight_db.store(keylight_metadatas.as_slice())
    }

    fn deduplicate_keylights(&mut self) {
        self.lights.sort_by_key(|light| light.metadata.ip.clone());
        self.lights.dedup_by_key(|light| light.metadata.ip.clone());
    }
}

#[cfg(test)]
mod test {

    use std::cell::RefCell;

    use crate::keylight::keylight_mocks::{
        create_metadata_list_fixture, MockKeylightAdapter, MockKeylightDb, MockKeylightFinder,
    };
    use crate::keylight::model::KeylightMetadata;

    use super::*;

    fn prepare_test() -> (MockKeylightFinder, MockKeylightAdapter, MockKeylightDb) {
        (
            MockKeylightFinder::new(create_metadata_list_fixture()),
            MockKeylightAdapter::new(Ok(vec![]), None),
            MockKeylightDb {
                stored_metadata_passed: RefCell::new(vec![]),
                load_response: Ok(create_metadata_list_fixture()),
            },
        )
    }

    mod discover_lights {
        use super::*;

        #[test]
        fn test_discover_keylights() {
            let (finder, adapter, db) = prepare_test();
            let deduplicated_metadata = vec![&finder.metadata[0], &finder.metadata[1]];
            let mut keylight_control = KeylightControl::new(&finder, &adapter, &db);
            keylight_control.discover_keylights();
            let discovered_metadata: Vec<&KeylightMetadata> = keylight_control
                .lights
                .iter()
                .map(|light| &light.metadata)
                .collect();
            assert_eq!(keylight_control.lights.len(), 2);
            assert_eq!(discovered_metadata, deduplicated_metadata);
        }
    }

    mod load_lights {
        use super::*;

        #[test]
        fn should_load_metadata_from_db() {
            let (finder, adapter, db) = prepare_test();
            let keylight_metadatas = db.load_response.clone().unwrap();
            let mut keylight_control = KeylightControl::new(&finder, &adapter, &db);
            let result = keylight_control.load_keylights();

            assert_eq!(result.is_ok(), true);
            for (index, keylight) in keylight_control.lights.iter().enumerate() {
                assert_eq!(keylight.metadata, keylight_metadatas[index]);
            }
        }

        #[test]
        fn should_discover_keylights_when_loading_from_db_fails() {
            let (finder, adapter, mut db) = prepare_test();
            db.load_response = Err(KeylightError::DbError(String::from("Test")));
            let mut keylight_control = KeylightControl::new(&finder, &adapter, &db);
            let result = keylight_control.load_keylights();

            assert_eq!(result.is_ok(), true);
            assert_eq!(keylight_control.lights.is_empty(), false);
            for (index, keylight) in keylight_control.lights.iter().enumerate() {
                assert_eq!(keylight.metadata, finder.metadata[index]);
            }
        }
    }
}
