import { Keylight as KeylightModel } from '../adapter';
import styles from './Keylight.module.css';

type KeyValueProps = {
  label: string;
  value: string;
};

type KeylightProps = {
  light: KeylightModel;
};

const KeyValue = ({ label, value }: KeyValueProps) => (
  <div class={styles.keyValue}>
    <label>{label}:</label>
    <span>{value}</span>
  </div>
);

export const Keylight = ({ light }: KeylightProps) => (
  <div class={styles.keylight}>
    <div class={styles.lightSwitch}>
      <button>{light.lights[0].on ? 'On' : 'Off'}</button>
    </div>
    <div class={styles.metadata}>
      <KeyValue label="Temperature" value={light.lights[0].temperature.toString()} />
      <KeyValue label="Brightness" value={light.lights[0].brightness.toString()} />
    </div>
  </div>
);
