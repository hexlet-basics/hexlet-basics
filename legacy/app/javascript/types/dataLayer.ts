/* eslint-disable @typescript-eslint/no-namespace */

export namespace DataLayer {
  /**
   * Base shape for any dataLayer event
   */
  export interface BaseEvent {
    event: string;
    [key: string]: unknown;
  }

  /**
   * eCommerce contract (GA4-style, suitable for most dataLayer consumers)
   */
  export namespace Ecommerce {
    export type Currency = string; // narrow if you want: 'USD' | 'EUR' | 'RUB'

    export interface Item {
      item_id: string;
      item_name: string;

      price: number;
      quantity: number;

      item_brand?: string;
      item_variant?: string;

      item_category?: string;
      item_category2?: string;
      item_category3?: string;
      item_category4?: string;
      item_category5?: string;

      discount?: number;
      coupon?: string;

      index?: number;
      item_list_id?: string;
      item_list_name?: string;
    }

    export interface Base {
      items: Item[];
    }

    export interface ViewItem extends Base {}
    export interface ViewItemList {
      item_list_id?: string;
      item_list_name?: string;
      items: Item[];
    }
    export interface SelectItem {
      item_list_id?: string;
      item_list_name?: string;
      items: Item[];
    }

    export interface AddToCart extends Base {}
    export interface RemoveFromCart extends Base {}
    export interface ViewCart extends Base {}

    export interface BeginCheckout extends Base {
      value?: number;
      currency?: Currency;
      coupon?: string;
    }

    export interface AddPaymentInfo extends Base {
      value?: number;
      currency?: Currency;
      coupon?: string;
      payment_type?: string;
    }

    export interface AddShippingInfo extends Base {
      value?: number;
      currency?: Currency;
      coupon?: string;
      shipping_tier?: string;
    }

    export interface Purchase extends Base {
      transaction_id: string;
      value: number;
      currency: Currency;

      affiliation?: string;
      tax?: number;
      shipping?: number;
      coupon?: string;
    }

    export interface Refund {
      transaction_id: string;
      value?: number;
      currency?: Currency;
      items?: Item[];
    }

    export type EventName =
      | "view_item"
      | "view_item_list"
      | "select_item"
      | "add_to_cart"
      | "remove_from_cart"
      | "view_cart"
      | "begin_checkout"
      | "add_payment_info"
      | "add_shipping_info"
      | "purchase"
      | "refund";

    export interface EventMap {
      view_item: ViewItem;
      view_item_list: ViewItemList;
      select_item: SelectItem;
      add_to_cart: AddToCart;
      remove_from_cart: RemoveFromCart;
      view_cart: ViewCart;
      begin_checkout: BeginCheckout;
      add_payment_info: AddPaymentInfo;
      add_shipping_info: AddShippingInfo;
      purchase: Purchase;
      refund: Refund;
    }

    export type Event<K extends EventName = EventName> = {
      event: K;
      ecommerce: EventMap[K];
    };
  }

  /**
   * Union type for what you allow to push into dataLayer
   * (eCommerce events + any other custom events)
   */
  export type Event = Ecommerce.Event | BaseEvent;
}
