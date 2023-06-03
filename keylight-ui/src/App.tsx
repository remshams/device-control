import { createSignal } from 'solid-js';
import './App.css';
import { Keylight, loadKeylights, refresh_lights } from './keylight/adapter';
import { KeylightList } from './keylight/components/KeylightList';
import { areLightsEditedSignal } from './keylight/stores';

const AppState = {
  init: 'init',
  loading: 'loading',
  loaded: 'loaded',
  error: 'error'
} as const;
type AppState = keyof typeof AppState;

export const Loading = () => <span>Loading</span>;

function App() {
  const [areLightsEdited] = areLightsEditedSignal;
  const [appState, setAppState] = createSignal<AppState>(AppState.loading);
  const [lights, setLights] = createSignal<Array<Keylight>>([]);
  loadKeylights().then(lights => {
    setLights(lights);
    setAppState(AppState.loaded);
    window.setInterval(() => {
      refresh_lights().then(lights => {
        if (!areLightsEdited()) {
          setLights(lights);
        }
      });
    }, 2000);
  });

  return <main>{appState() === AppState.loading ? <Loading /> : <KeylightList lights={lights()} />}</main>;
}

export default App;
