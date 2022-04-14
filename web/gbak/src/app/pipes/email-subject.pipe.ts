import {Pipe, PipeTransform} from "@angular/core";
import {Email} from "../definitions/email";
import {EmailHeaderPipe} from "./email-header.pipe";

@Pipe({name: 'emailSubject'})
export class EmailSubjectPipe extends EmailHeaderPipe {
  override header = 'Subject';
}
