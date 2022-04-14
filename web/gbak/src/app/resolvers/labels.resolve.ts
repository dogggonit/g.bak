import {Injectable} from "@angular/core";
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from "@angular/router";
import {Email} from "../definitions/email";
import {EmailService} from "../services/email.service";
import {catchError, EMPTY, map, Observable, tap} from "rxjs";
import {LabelsService} from "../services/labels.service";
import {APIUser} from "../definitions/api-email";

@Injectable()
export class LabelsResolver implements Resolve<APIUser[]> {
  constructor(
    private readonly labelService: LabelsService
  ) {}

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<APIUser[]> {
    return this.labelService.getLabels().pipe(
      catchError(() => EMPTY),
    );
  }
}
