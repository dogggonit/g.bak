import {Component, Input, OnInit} from "@angular/core";
import {ActivatedRoute, Router} from "@angular/router";
import {APIUser} from "../../definitions/api-email";

@Component({
  selector: 'app-toolbar',
  templateUrl: 'toolbar.component.html',
  styleUrls: [
    '../../app.component.css',
    'toolbar.component.css',
  ]
})
export class ToolbarComponent implements OnInit {
  labels: APIUser[] = [];
  sidenavOpened = false;

  @Input() title = '';
  @Input() subtitle = '';

  constructor(
    private readonly router: Router,
    private readonly route: ActivatedRoute,
  ) {
  }

  ngOnInit(): void {
    this.labels = this.route.snapshot.data['labels'];
  }

  handleSidenav = () => this.sidenavOpened = !this.sidenavOpened;
}
