import { invoke } from '@tauri-apps/api';

export type KeylightMetadata = {
  id: string;
  name: string;
  ip: number;
  port: number;
};

export type Keylight = {
  metadata: KeylightMetadata;
  lights: Array<Light>;
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
  const result = await invoke<Array<Keylight>>('discover_keylights').catch(_e => {
    throw KeylightError.keylightLoadError;
  });
  return result;
};

export const refresh_lights = async (): Promise<Array<Keylight>> => {
  const result = await invoke<Array<Keylight>>('refresh_lights').catch(_e => {
    throw KeylightError.keylightLoadError;
  });
  return result;
};

export const setLight = async (lightCommand: LightCommand): Promise<void> => {
  const result = await invoke<void>('set_light', { command: lightCommand }).catch(_e => {
    throw KeylightError.commandError;
  });
  return result;
};
