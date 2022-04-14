import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {Observable, of, switchMap, tap} from "rxjs";

@Injectable()
export class AddEmailService {
  url = "http://localhost:8080/api/v1/token";

  constructor (
    private readonly http: HttpClient,
  ) {}

  getTokenLink(): Observable<TokenRequest> {
    return this.http.get<TokenRequest>(this.url, {
      responseType: 'json',
      observe: 'body',
    });
  }

  postToken(code: string, token: string): Observable<void> {
    if (code === '' || token === '') {
      throw new Error('arguments cannot be blank');
    }

    return of({}).pipe(
      switchMap(() => this.http.post<void>(this.url, {
        Code: code,
        Token: token,
      })),
    );
  }
}

export interface TokenRequest {
  Code: string;
  URL: string;
}
