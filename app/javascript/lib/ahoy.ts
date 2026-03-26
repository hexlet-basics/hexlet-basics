import { noop } from "es-toolkit";
import { gon } from "./gon";
import { fromWindow } from "./utils";

type Ahoy = {
  configure(options: Record<string, unknown>): void;
  start(): void;
  trackView(options?: Record<string, unknown>): void;
};

const ahoyStub: Ahoy = {
  configure: noop,
  start: noop,
  trackView: noop,
};

let ahoyLoader: Promise<Ahoy> | null = null;

async function loadAhoy() {
  if (typeof window === "undefined") {
    return ahoyStub;
  }

  if (ahoyLoader) {
    return ahoyLoader;
  }

  ahoyLoader = import("ahoy.js").then(({ default: ahoy }) => {
    const realAhoy = ahoy as Ahoy;

    realAhoy.configure({
      startOnReady: false,
    });

    const ym = fromWindow("ym");
    if (ym) {
      ym(gon.ym_counter, "getClientID", (id: string) => {
        realAhoy.configure({
          visitParams: {
            ym_client_id: id,
          },
        });
        realAhoy.start();
      });
    }

    return realAhoy;
  });

  return ahoyLoader;
}

const ahoy: Ahoy = {
  configure(options) {
    void loadAhoy().then((instance) => instance.configure(options));
  },
  start() {
    void loadAhoy().then((instance) => instance.start());
  },
  trackView(options) {
    void loadAhoy().then((instance) => instance.trackView(options));
  },
};

export default ahoy;
