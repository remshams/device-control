import { Keylight } from '../adapter';
import { Keylight as KeylightComponent } from './Keylight';

export type KeylightListProps = {
  lights: Array<Keylight>;
};

export const KeylightList = (props: KeylightListProps) => {
  return (
    <>
      {props.lights.map(light => (
        <KeylightComponent light={light} />
      ))}
    </>
  );
};
