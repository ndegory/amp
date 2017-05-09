import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Organization } from '../models/organization.model';
import { MenuService } from '../services/menu.service';
import { OrganizationsService } from '../services/organizations.service';
import { ListService } from '../services/list.service';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../services/http.service';

@Component({
  selector: 'app-organizations',
  templateUrl: './organizations.component.html',
  styleUrls: ['./organizations.component.css'],
  providers: [ ListService ]
})
export class OrganizationsComponent implements OnInit {
  noOrganization = new Organization("", "")
  organization : Organization = this.noOrganization


  constructor(
    private route : ActivatedRoute,
    public organizationsService : OrganizationsService,
    public listService : ListService,
    private menuService : MenuService,
    private httpService : HttpService) {
      listService.setFilterFunction(organizationsService.match)
    }

  ngOnInit() {
    this.menuService.setItemMenu('organizations', 'List')
    this.listService.setData(this.organizationsService.organizations)
  }

  removeOrganization() {
    if(confirm("Are you sure to delete the organization: "+this.organization.name)) {
      this.menuService.waitingCursor(true)
      this.httpService.deleteOrganization(this.organization).subscribe(
        () => {
          this.menuService.waitingCursor(false)
          let list : Organization[] = []
          for (let org of this.organizationsService.organizations) {
            if (org.name != this.organization.name) {
              list.push(org)
            }
          }
          this.organizationsService.organizations=list
          this.listService.setData(this.organizationsService.organizations)
        },
        (error) => {
          this.menuService.waitingCursor(false)
          console.log(error)
        }
      )
    }
  }

  selectOrganization(org : Organization) {
    this.organization = org
  }

  createOrganization() {
    this.menuService.navigate(['/amp', 'organizations', 'create'])
  }

  editOrganization(org : Organization) {
    this.menuService.navigate(['/amp', 'organizations', org.name])
  }

}