import { Keylight as KeylightModel, setLight as setLightSystem } from '../adapter';
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

export const Keylight = (props: KeylightProps) => {
  const toggleOn = () => {
    setLightSystem({ id: props.light.metadata.id, index: 0, on: !props.light.lights[0].on });
  };
  return (
    <div class={styles.keylight}>
      <div class={styles.lightSwitch}>
        <button onClick={toggleOn}>{props.light.light.on ? 'On' : 'Off'}</button>
      </div>
      <div class={styles.metadata}>
        <KeyValue label="Temperature" value={props.light.light.temperature.toString()} />
        <KeyValue label="Brightness" value={props.light.light.brightness.toString()} />
      </div>
    </div>
  );
};
