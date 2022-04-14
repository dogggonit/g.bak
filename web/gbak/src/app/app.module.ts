import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {MatTableModule} from "@angular/material/table";
import {EmailService} from "./services/email.service";
import {HttpClient, HttpClientModule, HttpHandler} from "@angular/common/http";
import {MatDialogModule} from "@angular/material/dialog";
import {EmailViewComponent} from "./components/email-view/email-view.component";
import {MatPaginatorModule} from "@angular/material/paginator";
import {EmailHeaderPipe} from "./pipes/email-header.pipe";
import {EmailSubjectPipe} from "./pipes/email-subject.pipe";
import {EmailToPipe} from "./pipes/email-to.pipe";
import {EmailFromPipe} from "./pipes/email-from.pipe";
import {EmailMimePipe} from "./pipes/email-mime.pipe";
import {EmailTextBodyPipe} from "./pipes/email-text-body.pipe";
import {EmailHtmlBodyPipe} from "./pipes/email-html-body.pipe";
import {MatExpansionModule} from "@angular/material/expansion";
import {MatCardModule} from "@angular/material/card";
import {CommonModule} from "@angular/common";
import {MatTabsModule} from "@angular/material/tabs";
import { AppRoutingModule } from './app-routing.module';
import {EmailListComponent} from "./components/email-list/email-list.component";
import {RouterModule} from "@angular/router";
import {EmailResolver} from "./resolvers/email.resolve";
import {MatSnackBarModule} from "@angular/material/snack-bar";
import {MatToolbarModule} from "@angular/material/toolbar";
import {ToolbarComponent} from "./components/toolbar/toolbar.component";
import {MatButtonModule} from "@angular/material/button";
import {MatIconModule} from "@angular/material/icon";
import {MatSidenavModule} from "@angular/material/sidenav";
import {MatChipsModule} from "@angular/material/chips";
import {MatMenuModule} from "@angular/material/menu";
import {MatListModule} from "@angular/material/list";
import {FilesizePipe} from "./pipes/filesize.pipe";
import {Base64Pipe} from "./pipes/base64.pipe";
import {MatProgressSpinnerModule} from "@angular/material/progress-spinner";
import {HomeComponent} from "./components/home/home.component";
import {EmailListService} from "./services/email-list.service";
import {LabelsResolver} from "./resolvers/labels.resolve";
import {LabelsService} from "./services/labels.service";
import {ScrollingModule} from "@angular/cdk/scrolling";
import {AddEmailComponent} from "./components/add-email/add-email.component";
import {AddEmailService} from "./services/add-email.service";
import {MatInputModule} from "@angular/material/input";
import {ReactiveFormsModule} from "@angular/forms";
import {ToolbarService} from "./services/toolbar.service";
import {EmailPartBrowserComponent} from "./components/email-view/email-part-browser/email-part-browser.component";
import {EmailViewModule} from "./components/email-view/email-view.module";

@NgModule({
  declarations: [
    AppComponent,
    EmailSubjectPipe,
    EmailToPipe,
    EmailFromPipe,
    EmailMimePipe,
    EmailTextBodyPipe,
    EmailHtmlBodyPipe,
    EmailListComponent,
    ToolbarComponent,
    FilesizePipe,
    Base64Pipe,
    HomeComponent,
    AddEmailComponent,
  ],
  imports: [
    CommonModule,
    BrowserModule,
    BrowserAnimationsModule,
    MatTableModule,
    MatDialogModule,
    HttpClientModule,
    MatPaginatorModule,
    MatExpansionModule,
    MatCardModule,
    MatTabsModule,
    AppRoutingModule,
    RouterModule,
    MatSnackBarModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatSidenavModule,
    MatChipsModule,
    MatMenuModule,
    MatListModule,
    MatProgressSpinnerModule,
    ScrollingModule,
    MatInputModule,
    ReactiveFormsModule,
    EmailViewModule,
  ],
  providers: [
    ToolbarService,
    AddEmailService,
    EmailListService,
    LabelsService,
    EmailService,
    EmailHeaderPipe,
    EmailTextBodyPipe,
    EmailHtmlBodyPipe,
    EmailResolver,
    LabelsResolver,
    FilesizePipe,
    Base64Pipe,
    EmailSubjectPipe,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
