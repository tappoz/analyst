/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

const CompLibrary = require('../../core/CompLibrary.js');
const MarkdownBlock = CompLibrary.MarkdownBlock; /* Used to read markdown */
const Container = CompLibrary.Container;
const GridBlock = CompLibrary.GridBlock;

const siteConfig = require(process.cwd() + '/siteConfig.js');

function imgUrl(img) {
  return siteConfig.baseUrl + 'img/' + img;
}

function docUrl(doc, language) {
  return siteConfig.baseUrl + 'docs/' + (language ? language + '/' : '') + doc;
}

function pageUrl(page, language) {
  return siteConfig.baseUrl + (language ? language + '/' : '') + page;
}

class Button extends React.Component {
  render() {
    return (
      <div className="pluginWrapper buttonWrapper">
        <a className="button" href={this.props.href} target={this.props.target}>
          {this.props.children}
        </a>
      </div>
    );
  }
}

Button.defaultProps = {
  target: '_self',
};

const SplashContainer = props => (
  <div className="homeContainer">
    <div className="homeSplashFade">
      <div className="wrapper homeWrapper">{props.children}</div>
    </div>
  </div>
);

const Logo = props => (
  <div className="projectLogo">
    <img src={props.img_src} />
  </div>
);

const ProjectTitle = props => (
  <h2 className="projectTitle">
    {siteConfig.title}
    <small>{siteConfig.tagline}</small>
  </h2>
);

const PromoSection = props => (
  <div className="section promoSection">
    <div className="promoRow">
      <div className="pluginRowBlock">{props.children}</div>
    </div>
  </div>
);

class HomeSplash extends React.Component {
  render() {
    let language = this.props.language || '';
    return (
      <SplashContainer>
        
        <div className="inner">
          <ProjectTitle />
          <PromoSection>
            <Button href={docUrl('intro.html', language)}>Get Started</Button>
            <Button href={docUrl('intro.html', language)}>Examples</Button>
            <Button href={docUrl('blocks.html', language)}>Reference</Button>
          </PromoSection>
        </div>
      </SplashContainer>
    );
  }
}

const Block = props => (
  <Container
    padding={['bottom', 'top']}
    id={props.id}
    background={props.background}>
    <GridBlock align="center" contents={props.children} layout={props.layout} />
  </Container>
);

const Features = props => (
  <Block layout="fourColumn">
    {[
      {
        content: 'Like SQL, but easier to test and maintain. Write modular, easy to read code. Compile-time and execution-time testing features.',
        image: imgUrl('docusaurus.svg'),
        imageAlign: 'top',
        title: 'Easy to learn',
      },
      {
        content: 'Dozens of built-in connectors, including PostgreSQL, MySQL, MS SQL Server, HTTP APIs, Excel, and more.',
        image: imgUrl('docusaurus.svg'),
        imageAlign: 'top',
        title: 'Any data source',
      },
       {
        content: 'Use AQL for data movement, write control logic in any language. If you need to extend AQL, write a plugin in any language that supports JSON-RPC.',
        image: imgUrl('docusaurus.svg'),
        imageAlign: 'top',
        title: 'Extensible',
      },
      {
        content: 'Built-in transaction management that interrupts execution and rolls back changes if any issues are detected.',
        image: imgUrl('docusaurus.svg'),
        imageAlign: 'top',
        title: 'Safe by default',
      },
      {
        content: 'Log as much or as little information as you want, even down to every row of data that is processed. AQL admits Go templating syntax and parameters can be passed through command-line parameters, making AQL-based jobs easy to integrate into any scheduling or execution workflow.',
        image: imgUrl('docusaurus.svg'),
        imageAlign: 'top',
        title: 'Powerful Monitoring, super-easy scheduling',
      },
      {
        content: 'Download a single static binary and place it on your path. No messing around with database drivers, apt packages or makefiles',
        image: imgUrl('docusaurus.svg'),
        imageAlign: 'top',
        title: 'Easy installation on any platform',
      }
    ]}
  </Block>
);

const FeatureCallout = props => (
  <div
    className="productShowcaseSection paddingBottom"
    style={{textAlign: 'center'}}>

  </div>
);

const LearnHow = props => (
  <Block>
 
  </Block>
);

const TryOut = props => (
  <Block id="try">
  
  </Block>
);

const Description = props => (
  <Block>
  
  </Block>
);

const Showcase = props => {
  if ((siteConfig.users || []).length === 0) {
    return null;
  }
  const showcase = siteConfig.users
    .filter(user => {
      return user.pinned;
    })
    .map((user, i) => {
      return (
        <a href={user.infoLink} key={i}>
          <img src={user.image} title={user.caption} />
        </a>
      );
    });

  return (
    <div className="productShowcaseSection paddingBottom">
     
    </div>
  );
};

class Index extends React.Component {
  render() {
    let language = this.props.language || '';

    return (
      <div>
        <HomeSplash language={language} />
        <div className="mainContainer">
         
         
          <Showcase language={language} />
        </div>
      </div>
    );
  }
}

module.exports = Index;
