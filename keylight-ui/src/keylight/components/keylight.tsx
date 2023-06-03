import { createEffect, createSignal } from 'solid-js';
import { Keylight as KeylightModel, setKeylight } from '../adapter';
import { areLightsEditedSignal } from '../stores';
import styles from './Keylight.module.css';

const [_areLightsEdited, setAreLightsEdited] = areLightsEditedSignal;

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

type LightSliderProps = {
  label: string;
  value: number;
  min?: number;
  max?: number;
  onChange: (value: number) => void;
};

const LightSlider = (props: LightSliderProps) => (
  <div class={styles.lightSlider}>
    <label for="slider">{props.label}</label>
    <input
      id="slider"
      type="range"
      min={props.min ?? 0}
      max={props.max ?? 100}
      value={props.value}
      onChange={event => props.onChange(Number(event.currentTarget.value))}
      onMouseDown={_event => setAreLightsEdited(true)}
      onMouseUp={_event => setAreLightsEdited(false)}
    />
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
          <LightSlider
            label="Temperature"
            value={light().temperature}
            min={143}
            max={344}
            onChange={value => setLight({ ...light(), temperature: value })}
          />
          <KeyValue label="Temperature" value={light().temperature.toString()} />
        </div>
        <div>
          <LightSlider
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
