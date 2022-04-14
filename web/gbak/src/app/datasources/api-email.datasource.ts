import {CollectionViewer, DataSource} from "@angular/cdk/collections";
import {BehaviorSubject, catchError, finalize, Observable, of} from "rxjs";
import {EmailService} from "../services/email.service";
import {APIEmail, APIEmails} from "../definitions/api-email";
import {EmailListService} from "../services/email-list.service";


export class ApiEmailDatasource implements DataSource<APIEmail> {
  private emailSubject = new BehaviorSubject<APIEmail[]>([]);
  private lengthSubject = new BehaviorSubject<number>(0);
  private loadingSubject = new BehaviorSubject<boolean>(false);

  constructor(
    private emailListService: EmailListService,
    private emailAddress = '',
  ) {}

  connect(collectionViewer: CollectionViewer): Observable<APIEmail[]> {
    return this.emailSubject.asObservable();
  }

  disconnect(collectionViewer: CollectionViewer): void {
    this.emailSubject.complete();
    this.loadingSubject.complete();
  }

  loading(): Observable<boolean> {
    return this.loadingSubject.asObservable();
  }

  length(): Observable<number> {
    return this.lengthSubject.asObservable();
  }

  loadEmails(filter = '', sortDirection = 'desc', page = 1, limit = 25, label = 'INBOX') {
    this.loadingSubject.next(true);
    this.emailListService.getEmails(this.emailAddress, page, limit, label).pipe(
      catchError(() => of({Total: 0, Emails: []} as APIEmails)),
      finalize(() => this.loadingSubject.next(false))
    ).subscribe(emails => {
      this.lengthSubject.next(emails.Total);
      this.emailSubject.next(emails.Emails);
    });
  }
}
