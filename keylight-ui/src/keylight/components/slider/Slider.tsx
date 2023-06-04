import styles from './Slider.module.css';

export type SliderProps = {
  label: string;
  value: number;
  min?: number;
  max?: number;
  onChange: (value: number) => void;
  onChangeStart?: (value: number) => void;
  onChangeEnd?: (value: number) => void;
};

export const Slider = (props: SliderProps) => {
  const onchangeStart = props.onChangeStart ?? (() => {});
  const onChangeEnd = props.onChangeEnd ?? (() => {});
  return (
    <div class={styles.container}>
      <label for="slider">{props.label}:</label>
      <input
        id="slider"
        type="range"
        min={props.min ?? 0}
        max={props.max ?? 100}
        value={props.value}
        onChange={event => props.onChange(Number(event.currentTarget.value))}
        onMouseDown={event => onchangeStart(Number(event.currentTarget.value))}
        onMouseUp={event => onChangeEnd(Number(event.currentTarget.value))}
      />
    </div>
  );
};
