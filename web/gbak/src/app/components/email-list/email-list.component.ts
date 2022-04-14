import {AfterViewInit, Component, OnInit, ViewChild} from "@angular/core";
import {MatPaginator, PageEvent} from "@angular/material/paginator";
import {Email, Message} from "../../definitions/email";
import {EmailService} from "../../services/email.service";
import {MatDialog} from "@angular/material/dialog";
import {EmailViewComponent} from "../email-view/email-view.component";
import {ActivatedRoute, Router} from "@angular/router";
import {ApiEmailDatasource} from "../../datasources/api-email.datasource";
import {tap} from "rxjs";
import {APIEmail} from "../../definitions/api-email";
import {EmailListService} from "../../services/email-list.service";
import {ToolbarService} from "../../services/toolbar.service";

@Component({
  selector: 'email-list',
  templateUrl: 'email-list.component.html',
  styleUrls: [
    '../../app.component.css',
    'email-list.component.css',
  ]
})
export class EmailListComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator!: MatPaginator;

  dataSource!: ApiEmailDatasource;
  displayedColumns = ['From', 'Attachment', 'Subject', /*'Labels',*/ 'Date'];
  pageSizes = [25, 50, 100];

  label!: string;

  email!: string;

  loading = false;

  constructor(
    private readonly emailListService: EmailListService,
    private readonly dialog: MatDialog,
    private readonly router: Router,
    private readonly route: ActivatedRoute,
    private readonly toolbarService: ToolbarService,
  ) {  }

  ngOnInit(): void {
    this.dataSource = new ApiEmailDatasource(this.emailListService, this.route.snapshot.params['email']);
    this.dataSource.loading().subscribe(l => setTimeout(() => this.loading = l));
  }

  ngAfterViewInit() {
    this.dataSource.length().subscribe(l => {
      this.paginator.length = l;
    });

    this.route.params.subscribe(p => {
      this.label = p['label'];
      this.toolbarService.setTitle(this.label);
      this.loadEmailPage({
        pageIndex: 0,
        pageSize: this.pageSizes[0],
        length: 0,
      });
    });
  }

  rowClicked(email: APIEmail) {
    this.router.navigate(['emails', this.route.snapshot.params['email'], 'message', email.ID]);
  }

  loadEmailPage = (e: PageEvent) => {
    this.dataSource.loadEmails(
      '',
      'desc',
      e.pageIndex + 1,
      e.pageSize,
      this.label,
    );
  }
}
