import {Pipe, PipeTransform} from "@angular/core";
import {Email} from "../definitions/email";
import {EmailHeaderPipe} from "./email-header.pipe";
import {EmailMimePipe} from "./email-mime.pipe";

@Pipe({name: 'emailHtmlBody'})
export class EmailHtmlBodyPipe extends EmailMimePipe {
  override mimeType = 'text/html';
}
