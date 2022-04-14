import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {catchError, map, Observable} from "rxjs";
import {MatSnackBar} from "@angular/material/snack-bar";
import {APIEmail, APIEmails, APIUser} from "../definitions/api-email";
import {Email} from "../definitions/email";

@Injectable()
export class LabelsService {
  constructor (
    private readonly http: HttpClient,
  ) {}

  getLabels(): Observable<APIUser[]> {
    return this.http.get<APIUser[]>('http://localhost:8080/api/v1/labels', {
      responseType: 'json',
      observe: 'body',
    }).pipe(
      map(e => e.map(u => new APIUser(u))),
    );
  }
}
