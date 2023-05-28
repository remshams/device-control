import { createSignal } from 'solid-js';
import './App.css';
import { Keylight, loadKeylights } from './keylight/adapter';
import { KeylightList } from './keylight/components/KeylightList';

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

  return <main>{appState() === AppState.loading ? <Loading /> : <KeylightList lights={lights()} />}</main>;
}

export default App;
