import {Pipe, PipeTransform} from "@angular/core";
import {Email} from "../definitions/email";
import {EmailHeaderPipe} from "./email-header.pipe";

@Pipe({name: 'emailTo'})
export class EmailToPipe extends EmailHeaderPipe {
  override header = 'To';
}
