import {Pipe, PipeTransform} from "@angular/core";
import {Email, Message} from "../definitions/email";

@Pipe({name: 'filesize'})
export class FilesizePipe implements PipeTransform {
  transform = (size: number): string => {
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let unit = units[0];

    const makeStr = () => `${size.toFixed(2)} ${unit}`;

    for (let i = 0; i < units.length; i++) {
      size = Math.floor(size)
      if (size < 1024) {
        return makeStr();
      }
      size /= 1024;
      unit = units[Math.min(i+1, units.length-1)];
    }
    return makeStr();
  }
}
