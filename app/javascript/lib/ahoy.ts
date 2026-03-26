import ahoy from "ahoy.js";
import { gon } from "./gon";
import { fromWindow } from "./utils";

ahoy.configure({
  startOnReady: false,
});

const ym = fromWindow("ym");
if (ym) {
  ym(gon.ym_counter, "getClientID", (id: string) => {
    ahoy.configure({
      visitParams: {
        ym_client_id: id,
      },
    });
    ahoy.start();
  });
}

export default ahoy;
