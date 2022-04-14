import {Message} from "./email";

export interface EmailViewDisplay {
  displayType: DisplayType;
  message: Message;
}

export enum DisplayTypeList {
  body = 'body',
  header = 'header',
  part = 'part',
  attachment = 'attachment',
}

export type DisplayType = keyof typeof DisplayTypeList;
