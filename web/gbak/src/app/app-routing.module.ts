import { NgModule } from '@angular/core';
import {RouterModule, Routes} from "@angular/router";
import {EmailListComponent} from "./components/email-list/email-list.component";
import {EmailViewComponent} from "./components/email-view/email-view.component";
import {EmailResolver} from "./resolvers/email.resolve";
import {HomeComponent} from "./components/home/home.component";
import {LabelsResolver} from "./resolvers/labels.resolve";

const routes: Routes = [
  { path: 'emails/:email/label/:label', component: EmailListComponent },
  { path: 'emails/:email/message/:id', component: EmailViewComponent, resolve: {email: EmailResolver} },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
