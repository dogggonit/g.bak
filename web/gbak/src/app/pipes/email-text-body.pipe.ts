import {Pipe, PipeTransform} from "@angular/core";
import {Email} from "../definitions/email";
import {EmailHeaderPipe} from "./email-header.pipe";
import {EmailMimePipe} from "./email-mime.pipe";

@Pipe({name: 'emailTextBody'})
export class EmailTextBodyPipe extends EmailMimePipe {
  override mimeType = 'text/plain';
}
