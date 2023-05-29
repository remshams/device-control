import { createEffect, createSignal } from 'solid-js';
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
  const [light, setLight] = createSignal(props.light.light);
  const toggleOn = () => {
    setLight({ ...light(), on: !light().on });
  };
  createEffect(() => {
    setLightSystem({ id: props.light.metadata.id, index: 0, ...light() });
  });

  return (
    <div class={styles.keylight}>
      <div class={styles.lightSwitch}>
        <button onClick={toggleOn}>{light().on ? 'On' : 'Off'}</button>
      </div>
      <div class={styles.metadata}>
        <KeyValue label="Temperature" value={light().temperature.toString()} />
        <KeyValue label="Brightness" value={light().brightness.toString()} />
      </div>
    </div>
  );
};
