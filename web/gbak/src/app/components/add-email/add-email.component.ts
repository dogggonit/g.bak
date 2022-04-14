import {Component, OnInit, ViewChild} from "@angular/core";
import {AddEmailService, TokenRequest} from "../../services/add-email.service";
import {MatDialogRef} from "@angular/material/dialog";
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {MatInput} from "@angular/material/input";
import {catchError, EMPTY} from "rxjs";

@Component({
  templateUrl: 'add-email.component.html',
  selector: 'app-add-email-dialog',
  styleUrls: [
    'add-email.component.css',
  ],
})
export class AddEmailComponent implements OnInit {
  options!: FormGroup;
  tokenInput = new FormControl('');

  req: TokenRequest = {URL: '', Code: ''};
  loading = false;

  constructor(
    private readonly addEmailService: AddEmailService,
    public dialogRef: MatDialogRef<AddEmailComponent>,
    fb: FormBuilder,
  ) {
    this.options = fb.group({
      token: this.tokenInput,
    });
  }

  ngOnInit(): void {
    this.addEmailService.getTokenLink().subscribe(t => this.req = t);
  }

  onOk = () => {
    this.loading = true;
    this.addEmailService.postToken(this.req.Code, this.tokenInput.value).pipe(
      catchError(() => {
        this.loading = false;
        return EMPTY;
      }),
    ).subscribe(() => {
      this.loading = false;
      this.dialogRef.close();
    });
  }

  onCancel = () => {
    this.dialogRef.close();
    this.loading = false;
  }

  openSignIn = () => {
    window.open(this.req.URL, '_blank');
  }
}
