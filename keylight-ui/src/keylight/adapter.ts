import { invoke } from '@tauri-apps/api';
import { convertKeylightDto, KeylightDto } from './converter';

export type KeylightMetadata = {
  id: string;
  name: string;
  ip: number;
  port: number;
};

export type Keylight = {
  metadata: KeylightMetadata;
  lights: Array<Light>;
  light: Light;
};

export type Light = {
  on: boolean;
  temperature: number;
  brightness: number;
};

export type LightCommand = {
  id: string;
  index: number;
  on?: boolean;
  temperature?: number;
  brightness?: number;
};

export const KeylightError = {
  keylightLoadError: 'keylightLoadError',
  commandError: 'commandError'
} as const;
export type KeylightError = keyof typeof KeylightError;

export const loadKeylights = async (): Promise<Array<Keylight>> => {
  const result = await invoke<Array<KeylightDto>>('discover_keylights').catch(_e => {
    throw KeylightError.keylightLoadError;
  });
  return result.map(convertKeylightDto);
};

export const refresh_lights = async (): Promise<Array<Keylight>> => {
  const result = await invoke<Array<KeylightDto>>('refresh_lights').catch(_e => {
    throw KeylightError.keylightLoadError;
  });
  return result.map(convertKeylightDto);
};

export const setKeylight = async (lightCommand: LightCommand): Promise<void> => {
  const result = await invoke<void>('set_light', { command: lightCommand }).catch(_e => {
    throw KeylightError.commandError;
  });
  return result;
};
