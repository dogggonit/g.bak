import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {map, Observable} from "rxjs";
import {APIEmails} from "../definitions/api-email";

@Injectable()
export class EmailListService {
  constructor (
    private readonly http: HttpClient,
  ) {}

  getEmails(email: string, page: number, limit: number, label: string): Observable<APIEmails> {
    return this.http.get<APIEmails>(`http://localhost:8080/api/v1/emails/${email}/label/${label}`, {
      params: {
        page,
        limit,
      },
      responseType: 'json',
      observe: 'body',
    }).pipe(
      map(ae => new APIEmails(ae)),
    );
  }
}
