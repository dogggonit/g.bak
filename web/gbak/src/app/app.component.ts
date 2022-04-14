import {Component, OnInit} from '@angular/core';
import {APILabel, APIUser} from "./definitions/api-email";
import {ActivatedRoute, NavigationEnd, NavigationStart, Router} from "@angular/router";
import {LabelsService} from "./services/labels.service";
import {MatDialog} from "@angular/material/dialog";
import {AddEmailComponent} from "./components/add-email/add-email.component";
import {ToolbarService} from "./services/toolbar.service";
import {BehaviorSubject} from "rxjs";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  labels: APIUser[] = [];
  title = '';
  opened = new BehaviorSubject(false);

  constructor(
    private readonly router: Router,
    private readonly labelsService: LabelsService,
    private readonly dialog: MatDialog,
    private readonly route: ActivatedRoute,
    private readonly toolbarService: ToolbarService,
  ) {
  }

  ngOnInit(): void {
    this.router.events.subscribe(e => {
      if (e instanceof NavigationStart) {
        this.opened.next(false);
      }

      if (e instanceof NavigationEnd) {
        this.labelsService.getLabels().subscribe(l => this.labels = l);
      }
    });

    this.route.data.subscribe(d => {
      if (typeof d['title'] === 'string') {
        this.title = d['title'];
      }
    });

    this.toolbarService.getTitle().subscribe(t => this.title = t);
  }

  handleSidenav = () => this.opened.next(!this.opened.value);

  addNew = () => {
    const dialogRef = this.dialog.open(AddEmailComponent);
  }

  labelClicked(u: APIUser, l: APILabel) {
    this.router.navigate(['emails', u.Email, 'label', l.GmailId])
  }
}
