import {
  Flex,
  Image,
  ScrollArea,
  Sidebar as SidebarComponent,
  Text,
  TextField
} from '@raystack/apsara';
import { useCallback, useState } from 'react';
import { Link, useRouterState } from '@tanstack/react-router';
import organization from '~/react/assets/organization.png';
import user from '~/react/assets/user.png';
import { organizationNavItems, userNavItems } from './helpers';

// @ts-ignore
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';
import styles from './sidebar.module.css';

export const Sidebar = () => {
  const [search, setSearch] = useState('');
  const routerState = useRouterState();

  const isActive = useCallback(
    (path: string) =>
      path.length > 2
        ? routerState.location.pathname.includes(path)
        : routerState.location.pathname === path,
    [routerState.location.pathname]
  );

  return (
    <SidebarComponent>
      <ScrollArea className={styles.scrollarea}>
        <Flex direction="column" style={{ gap: '24px', marginTop: '40px' }}>
          <TextField
            // @ts-ignore
            size="medium"
            // @ts-ignore
            leading={
              <MagnifyingGlassIcon
                style={{ color: 'var(--foreground-base)' }}
              />
            }
            placeholder="Search"
            onChange={event => setSearch(event.target.value)}
          />
          <SidebarComponent.Navigations>
            <SidebarComponent.NavigationGroup
              name="Organization"
              icon={
                <Image
                  alt="organization"
                  width={16}
                  height={16}
                  src={organization}
                />
              }
            >
              {organizationNavItems
                .filter(s => s.name.toLowerCase().includes(search))
                .map(nav => {
                  return (
                    <SidebarComponent.NavigationCell
                      key={nav.name}
                      asChild
                      active={!!isActive(nav?.to as string) as any}
                      style={{ padding: 0 }}
                    >
                      <Link
                        key={nav.name}
                        to={nav.to as string}
                        style={{
                          width: '100%',
                          textDecoration: 'none',
                          padding: 'var(--pd-8)'
                        }}
                        search={{}}
                        params={{}}
                      >
                        <Text
                          style={{
                            color: 'var(--foreground-base)',
                            fontWeight: '500'
                          }}
                        >
                          {nav.name}
                        </Text>
                      </Link>
                    </SidebarComponent.NavigationCell>
                  );
                })}
            </SidebarComponent.NavigationGroup>
            <SidebarComponent.NavigationGroup
              name="My Account"
              icon={<Image alt="user" width={16} height={16} src={user} />}
            >
              {userNavItems
                .filter(s => s.name.toLowerCase().includes(search))
                .map(nav => (
                  <SidebarComponent.NavigationCell
                    key={nav.name}
                    asChild
                    active={!!isActive(nav?.to as string) as any}
                    style={{ padding: 0 }}
                  >
                    <Link
                      key={nav.name}
                      to={nav.to as string}
                      style={{
                        width: '100%',
                        textDecoration: 'none',
                        padding: 'var(--pd-8)'
                      }}
                      search={{}}
                      params={{}}
                    >
                      <Text
                        style={{
                          color: 'var(--foreground-base)',
                          fontWeight: '500'
                        }}
                      >
                        {nav.name}
                      </Text>
                    </Link>
                  </SidebarComponent.NavigationCell>
                ))}
            </SidebarComponent.NavigationGroup>
          </SidebarComponent.Navigations>
        </Flex>
      </ScrollArea>
    </SidebarComponent>
  );
};

type SidebarHeaderProps = { children?: React.ReactNode };
export const SidebarHeader = ({ children }: SidebarHeaderProps) => {
  return <Flex justify="between">{children}</Flex>;
};

type SidebarFooterProps = { children?: React.ReactNode };
export const SidebarFooter = ({ children }: SidebarFooterProps) => {
  return <Flex>{children}</Flex>;
};
