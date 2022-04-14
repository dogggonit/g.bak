import {Pipe, PipeTransform} from "@angular/core";
import {Email, Message} from "../definitions/email";

@Pipe({name: 'base64'})
export class Base64Pipe implements PipeTransform {
  transform = (base64: string): string => {
    base64 = base64.replace(/_/g, '/');
    base64 = base64.replace(/-/g, '+');
    return atob(base64.replace(/\s/g, ''));
  }
}
