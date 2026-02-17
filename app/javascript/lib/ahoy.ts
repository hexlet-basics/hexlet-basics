import ahoy from "ahoy.js";
import { fromWindow } from "./utils";

const ymClientId = fromWindow("ymClientId");
ahoy.configure({
  visitParams: {
    ym_client_id: ymClientId,
  },
});

export default ahoy;
