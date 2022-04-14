import {NgModule} from "@angular/core";
import {EmailViewComponent} from "./email-view.component";
import {EmailPartBrowserComponent} from "./email-part-browser/email-part-browser.component";
import {MatLineModule} from "@angular/material/core";
import {MatCardModule} from "@angular/material/card";
import {EmailSubjectPipe} from "../../pipes/email-subject.pipe";
import {EmailAttachmentListComponent} from "./email-part-browser/email-attachment-list/email-attachment-list.component";
import {MatButtonModule} from "@angular/material/button";
import {CommonModule} from "@angular/common";
import {EmailBodyListComponent} from "./email-part-browser/email-body-list/email-body-list.component";
import {EmailPartDisplayComponent} from "./email-part-display/email-part-display.component";


@NgModule({
  declarations: [
    EmailViewComponent,
    EmailPartBrowserComponent,
    EmailAttachmentListComponent,
    EmailBodyListComponent,
    EmailPartDisplayComponent,
  ],
  imports: [
    CommonModule,
    MatLineModule,
    MatCardModule,
    MatButtonModule
  ],
  providers: [
    EmailSubjectPipe,
  ],
  exports: [
    EmailViewComponent,
  ],
})
export class EmailViewModule {

}
