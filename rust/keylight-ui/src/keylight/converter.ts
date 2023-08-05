import { Keylight, KeylightMetadata, Light } from './adapter';

export type KeylightDto = {
  metadata: KeylightMetadata;
  lights: Array<Light>;
};

export const convertKeylightDto = (dto: KeylightDto): Keylight => {
  return {
    metadata: dto.metadata,
    light: dto.lights[0],
    lights: dto.lights
  };
};
