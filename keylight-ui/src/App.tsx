import { createSignal } from 'solid-js';
import './App.css';
import { Keylight, loadKeylights } from './keylight/adapter';
import { Keylight as KeylightComponent } from './keylight/components/keylight';

const AppState = {
  init: 'init',
  loading: 'loading',
  loaded: 'loaded',
  error: 'error'
} as const;
type AppState = keyof typeof AppState;

export const Loading = () => <span>Loading</span>;

function App() {
  const [appState, setAppState] = createSignal<AppState>(AppState.loading);
  const [lights, setLights] = createSignal<Array<Keylight>>([]);
  loadKeylights().then(lights => {
    setLights(lights);
    setAppState(AppState.loaded);
  });

  return <main>{appState() === AppState.loading ? <Loading /> : <KeylightComponent light={lights()[0]} />}</main>;
}

export default App;
