import {Pipe, PipeTransform} from "@angular/core";
import {Email, Message} from "../definitions/email";

@Pipe({name: 'emailHeader'})
export class EmailHeaderPipe implements PipeTransform {
  header = '';

  transform = (email: Email, header = this.header): string => {
    return email.Message ? this.getHeader(email.Message, header) : '';
  }

  getHeader = (msg: Message, name: string): string => {
    for (const h of msg.Headers) {
      if (h.Name === name) {
        return h.Value;
      }
    }

    for (const p of msg.Parts) {
      const h = this.getHeader(p, name);
      if (h != '') {
        return h;
      }
    }

    return '';
  }
}
