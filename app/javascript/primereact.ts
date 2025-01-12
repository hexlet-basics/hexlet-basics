import classNames from "classnames";
import type { PrimeReactPTOptions } from "primereact/api";

const TRANSITIONS = {
  overlay: {
    timeout: 150,
    classNames: {
      enter: "opacity-0 scale-75",
      enterActive:
        "opacity-100 scale-100 transition-transform transition-opacity duration-150 ease-in",
      exit: "opacity-100",
      exitActive: "opacity-0 transition-opacity duration-150 ease-linear",
    },
  },
};

// TODO: Extract to https://github.com/Hexlet/primereact-bootstrap-theme
const pt: PrimeReactPTOptions = {
  datatable: {
    table: { className: "table" },
  },
  button: {
    root: ({ props, context }) => ({
      className: classNames(
        "btn d-inline-flex align-items-center justify-center text-center",
        "transition duration-200 ease-in-out",
        {
          "btn-primary":
            !props.link &&
            props.severity === null &&
            !props.text &&
            !props.outlined &&
            !props.plain,
          "btn-secondary":
            props.severity === "secondary" &&
            !props.text &&
            !props.outlined &&
            !props.plain,
          "btn-success":
            props.severity === "success" &&
            !props.text &&
            !props.outlined &&
            !props.plain,
          "btn-info":
            props.severity === "info" &&
            !props.text &&
            !props.outlined &&
            !props.plain,
          "btn-warning":
            props.severity === "warning" &&
            !props.text &&
            !props.outlined &&
            !props.plain,
          "btn-danger":
            props.severity === "danger" &&
            !props.text &&
            !props.outlined &&
            !props.plain,
          "btn-outline-secondary":
            props.outlined && props.severity === "secondary",
          "btn-outline-success": props.outlined && props.severity === "success",
          "btn-outline-info": props.outlined && props.severity === "info",
          "btn-outline-warning": props.outlined && props.severity === "warning",
          "btn-outline-danger": props.outlined && props.severity === "danger",
          "btn-link": props.link,
          "btn-sm": props.size === "small",
          "btn-lg": props.size === "large",
          "btn-block": props.block,
          "disabled opacity-50": context.disabled,
        },
        {
          rounded: !props.rounded,
          "rounded-pill": props.rounded,
        },
        {
          "d-flex flex-column":
            props.iconPos === "top" || props.iconPos === "bottom",
        },
        { "pe-5": props.loading && props.label != null },
      ),
    }),
    label: ({ props }) => ({
      className: classNames(
        "flex-1 fw-bold",
        { "text-decoration-underline": props.link },
        { invisible: !props.label },
      ),
    }),
    icon: ({ props }) => ({
      className: classNames("me-2", {
        "ms-2": props.iconPos === "right" && props.label != null,
        "mb-2": props.iconPos === "top" && props.label != null,
        "mt-2": props.iconPos === "bottom" && props.label != null,
      }),
    }),
    loadingIcon: ({ props }) => ({
      className: classNames("spinner-border spinner-border-sm me-2", {
        "ms-2":
          props.loading && props.iconPos === "right" && props.label != null,
        "mb-2": props.loading && props.iconPos === "top" && props.label != null,
        "mt-2":
          props.loading && props.iconPos === "bottom" && props.label != null,
      }),
    }),
    badge: ({ props }) => ({
      className: classNames("badge rounded-pill ms-2", {
        "badge-primary": props.badge && props.severity === "primary",
        "badge-secondary": props.badge && props.severity === "secondary",
        "badge-success": props.badge && props.severity === "success",
        "badge-warning": props.badge && props.severity === "warning",
        "badge-danger": props.badge && props.severity === "danger",
      }),
    }),
  },

  dropdown: {
    root: ({ props }) => ({
      className: classNames(
        "dropdown position-relative d-inline-flex",
        "bg-white border border-secondary rounded",
        "transition-all ease-in-out",
        "w-100 md:w-56",
        "hover:border-primary focus:outline-none focus-shadow",
        { "opacity-50 pointer-events-none cursor-default": props.disabled },
      ),
    }),
    input: ({ props }) => ({
      className: classNames(
        "form-control border-0",
        "text-body bg-transparent",
        "py-2 px-3",
        "transition-colors bg-transparent rounded",
        { "pe-5": props.showClear },
      ),
    }),
    trigger: {
      className: classNames(
        "btn btn-outline-secondary d-flex align-items-center justify-content-center",
        "rounded-end",
      ),
    },
    wrapper: {
      className: classNames(
        "dropdown-menu shadow-lg",
        "max-h-200 overflow-auto",
        "bg-white text-body rounded",
      ),
    },
    list: "list-unstyled p-2",
    item: ({ context }) => ({
      className: classNames(
        "dropdown-item d-flex align-items-center",
        "transition-all",
        {
          "text-body": !context.focused && !context.selected,
          "active bg-primary text-white": context.focused && context.selected,
          "bg-light": context.focused && !context.selected,
          "disabled opacity-50": context.disabled,
        },
      ),
    }),
    itemgroup: {
      className: classNames("dropdown-header fw-bold", "text-muted"),
    },
    header: {
      className: classNames(
        "dropdown-header text-muted bg-light",
        "rounded-top",
      ),
    },
    filtercontainer: "position-relative",
    filterinput: {
      className: classNames(
        "form-control",
        "w-100",
        "py-2 px-3",
        "rounded",
        "border-secondary",
        "transition-all",
      ),
    },
    filtericon: "position-absolute top-50 translate-middle-y",
    clearicon: "position-absolute top-50 end-0 translate-middle-y text-muted",
    transition: TRANSITIONS.overlay,
  },
  autocomplete: {
    root: ({ props }) => ({
      className: classNames(
        "position-relative d-inline-flex",
        {
          "opacity-50 pointer-events-none": props.disabled,
        },
        { "w-100": props.multiple },
      ),
    }),
    container: {
      className: classNames(
        "list-unstyled d-flex align-items-center flex-wrap w-100",
        "px-3 py-2 gap-2",
        "text-body bg-white border border-secondary rounded",
        "focus-within:shadow focus-within:border-primary",
        "transition-all",
      ),
    },
    inputToken: {
      className: classNames("py-1 px-2", "flex-grow-1 d-inline-flex"),
    },
    input: ({ props }) => ({
      root: {
        className: classNames(
          "m-0",
          "rounded",
          { "rounded-end-0": props.dropdown },
          {
            "form-control": !props.multiple,
            "form-control-plaintext m-0 p-0 w-100": props.multiple,
          },
        ),
      },
    }),
    token: {
      className: classNames(
        "py-1 px-2 me-2 bg-secondary text-white rounded-pill",
        "d-inline-flex align-items-center",
      ),
    },
    dropdownButton: {
      root: "rounded-start-0",
    },
    panel: {
      className: classNames(
        "bg-white text-body border-0 rounded shadow",
        "overflow-auto",
        "max-h-200px",
      ),
    },
    list: "py-2 list-unstyled m-0",
    item: ({ context }) => ({
      className: classNames(
        "cursor-pointer text-body overflow-hidden",
        "py-2 px-3 transition-shadow",
        {
          "text-muted": !context.selected,
          "bg-primary text-white": context.selected,
          "bg-light": context.focused && !context.selected,
        },
      ),
    }),
    itemGroup: {
      className: classNames("py-2 px-3 fw-bold text-secondary bg-light"),
    },
    transition: TRANSITIONS.overlay,
  },
};

export default pt;
