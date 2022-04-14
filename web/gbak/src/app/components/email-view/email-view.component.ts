import {Component, Inject, OnInit} from "@angular/core";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {Email, Message} from "../../definitions/email";
import {EmailTextBodyPipe} from "../../pipes/email-text-body.pipe";
import {EmailHtmlBodyPipe} from "../../pipes/email-html-body.pipe";
import {ActivatedRoute, Router} from "@angular/router";
import {ToolbarService} from "../../services/toolbar.service";
import {EmailSubjectPipe} from "../../pipes/email-subject.pipe";
import {EmailViewDisplay} from "../../definitions/internal";

@Component({
  selector: 'email-view',
  templateUrl: 'email-view.component.html',
  styleUrls: [
    '../../app.component.css',
    'email-view.component.css',
  ]
})
export class EmailViewComponent implements OnInit {
  email!: Email;
  toDisplay?: EmailViewDisplay;

  constructor(
    private readonly emailSubjectPipe: EmailSubjectPipe,
    private readonly route: ActivatedRoute,
    private readonly toolbarService: ToolbarService,
  ) {
  }

  ngOnInit(): void {
    this.email = this.route.snapshot.data['email'];
    this.toolbarService.setTitle(this.emailSubjectPipe.transform(this.email));
  }

  get message(): Message {
    return this.email.Message || {} as Message;
  }
}
