import { Keylight } from '../adapter';
import { Keylight as KeylightComponent } from './Keylight';

export type KeylightListProps = {
  lights: Array<Keylight>;
};

export const KeylightList = ({ lights }: KeylightListProps) => (
  <>
    {lights.map(light => (
      <KeylightComponent light={light} />
    ))}
  </>
);
