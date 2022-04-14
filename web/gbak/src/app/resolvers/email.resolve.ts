import {Injectable} from "@angular/core";
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from "@angular/router";
import {Email} from "../definitions/email";
import {EmailService} from "../services/email.service";
import {catchError, EMPTY, map, Observable, tap} from "rxjs";

@Injectable()
export class EmailResolver implements Resolve<Email> {
  constructor(
    private readonly emailService: EmailService
  ) {}

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Email> {
    return this.emailService.getEmail(route.params['email'], Number(route.params['id'])).pipe(
      catchError(() => EMPTY),
    );
  }
}
