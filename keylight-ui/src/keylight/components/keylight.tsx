import { createEffect, createSignal } from 'solid-js';
import { Keylight as KeylightModel, setKeylight } from '../adapter';
import { areLightsEditedSignal } from '../stores';
import styles from './Keylight.module.css';
import { Slider } from './slider/Slider';

const [_areLightsEdited] = areLightsEditedSignal;

type KeyValueProps = {
  label: string;
  value: string;
};

type KeylightProps = {
  light: KeylightModel;
};

const KeyValue = (props: KeyValueProps) => (
  <div class={styles.keyValue}>
    <label>{props.label}:</label>
    <span>{props.value}</span>
  </div>
);

export const Keylight = (props: KeylightProps) => {
  const [light, setLight] = createSignal(props.light.light);
  createEffect(() => {
    setKeylight({ id: props.light.metadata.id, index: 0, ...light() });
  });

  return (
    <div class={styles.keylight}>
      <div class={styles.lightSwitch}>
        <button onClick={() => setLight({ ...light(), on: !light().on })}>{light().on ? 'On' : 'Off'}</button>
      </div>
      <div class={styles.metadata}>
        <div>
          <Slider
            label="Temperature"
            value={light().temperature}
            min={143}
            max={344}
            onChange={value => setLight({ ...light(), temperature: value })}
          />
          <KeyValue label="Temperature" value={light().temperature.toString()} />
        </div>
        <div>
          <Slider
            label="Brightness"
            value={light().brightness}
            onChange={value => setLight({ ...light(), brightness: value })}
          />
          <KeyValue label="Brightness" value={light().brightness.toString()} />
        </div>
      </div>
    </div>
  );
};
