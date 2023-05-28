import { Keylight as KeylightModel } from '../adapter';
import styles from './Keylight.module.css';

export const Keylight = ({ light }: { light: KeylightModel }) => {
  return (
    <div class={styles.keylight}>
      <div class={styles.lightSwitch}>
        <button>{light.lights[0].on ? 'On' : 'Off'}</button>
      </div>
      <div class={styles.metadata}>
        <div>
          <label>Temperature:</label>
          <span>{light.lights[0].temperature}</span>
        </div>
        <div>
          <label>Brightness:</label>
          <span>{light.lights[0].brightness}</span>
        </div>
      </div>
    </div>
  );
};
