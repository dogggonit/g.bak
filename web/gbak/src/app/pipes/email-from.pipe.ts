import {Pipe, PipeTransform} from "@angular/core";
import {Email} from "../definitions/email";
import {EmailHeaderPipe} from "./email-header.pipe";

@Pipe({name: 'emailFrom'})
export class EmailFromPipe extends EmailHeaderPipe {
  override header = 'From';
}
