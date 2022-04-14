import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {catchError, map, Observable} from "rxjs";
import {MatSnackBar} from "@angular/material/snack-bar";
import {APIEmail, APIEmails} from "../definitions/api-email";
import {Email} from "../definitions/email";

@Injectable()
export class EmailService {
  constructor (
    private readonly http: HttpClient,
  ) {}

  getEmail(email: string, id: number): Observable<Email> {
    return this.http.get<Email>(`http://localhost:8080/api/v1/emails/${email}/message/${id}`, {
      responseType: 'json',
      observe: 'body',
    }).pipe(
      map(e => new Email(e)),
    );
  }
}
